package field

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var ErrUnsupportedKind = errors.New("unsupported kind")

var defaultBinder = NewBinder()

func Bind(ctx context.Context, message proto.Message) error {
	return defaultBinder.Bind(ctx, message)
}

func BindHeader(header http.Header, message proto.Message) error {
	return defaultBinder.BindHeader(header, message)
}

func BindMetadata(md metadata.MD, message proto.Message) error {
	return defaultBinder.BindMetadata(md, message)
}

type headerRule struct {
	fd     protoreflect.FieldDescriptor
	header string
}

type Binder struct {
	enableCache bool
	cache       *sync.Map
}
type BindOption func(b *Binder)

func WithCache(cache bool) BindOption {
	return func(b *Binder) {
		b.enableCache = cache
	}
}

func NewBinder(options ...BindOption) *Binder {
	b := &Binder{enableCache: false}
	for _, option := range options {
		option(b)
	}
	if b.enableCache {
		b.cache = new(sync.Map)
	}
	return b
}

func (b *Binder) Bind(ctx context.Context, message proto.Message) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return b.bind(message, func(key string) []string { return md.Get(key) })
}

func (b *Binder) BindHeader(header http.Header, message proto.Message) error {
	if len(header) == 0 {
		return nil
	}
	return b.bind(message, func(key string) []string { return header.Values(key) })
}

func (b *Binder) BindMetadata(md metadata.MD, message proto.Message) error {
	if len(md) == 0 {
		return nil
	}
	return b.bind(message, func(key string) []string { return md.Get(key) })
}

func (b *Binder) bind(message proto.Message, kv func(string) []string) error {
	m := message.ProtoReflect()
	descriptor := m.Descriptor()

	var rules []*headerRule
	if b.enableCache {
		if val, ok := b.cache.Load(descriptor.FullName()); ok {
			rules = val.([]*headerRule)
		} else {
			rules = b.decodeFields(descriptor)
			b.cache.Store(descriptor.FullName(), rules)
		}
	} else {
		rules = b.decodeFields(descriptor)
	}

	for _, rule := range rules {
		if err := b.setHeader(m, rule.fd, kv(rule.header)); err != nil {
			return err
		}
	}
	return nil
}

func (b *Binder) decodeFields(descriptor protoreflect.MessageDescriptor) []*headerRule {
	var rules []*headerRule
	for i := 0; i < descriptor.Fields().Len(); i++ {
		field := descriptor.Fields().Get(i)
		opts := field.Options().(*descriptorpb.FieldOptions)
		if !proto.HasExtension(opts, E_Binding) {
			continue
		}
		binding, ok := proto.GetExtension(opts, E_Binding).(*Binding)
		if !ok {
			continue
		}
		rules = append(rules, &headerRule{
			fd:     field,
			header: binding.Header,
		})
	}
	return rules
}

func (b *Binder) setHeader(m protoreflect.Message, fd protoreflect.FieldDescriptor, values []string) error {
	if len(values) == 0 {
		return nil //skip empty value
	}
	if fd.IsList() {
		list := m.Mutable(fd).List()
		for _, raw := range values {
			v, err := b.parse(fd.Kind(), raw)
			if errors.Is(err, ErrUnsupportedKind) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("field %s: %w", fd.FullName(), err)
			}
			list.Append(v)
		}
	} else {
		v, err := b.parse(fd.Kind(), values[0])
		if errors.Is(err, ErrUnsupportedKind) {
			return nil
		}
		if err != nil {
			return err
		}
		m.Set(fd, v)

	}
	return nil
}

func (b *Binder) parse(kind protoreflect.Kind, value string) (protoreflect.Value, error) {
	switch kind {
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(value), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt32(int32(v)), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt64(v), nil
	case protoreflect.BoolKind:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfBool(v), nil
	default:
		return protoreflect.Value{}, fmt.Errorf("%w: %s", ErrUnsupportedKind, kind.String())
	}
}
