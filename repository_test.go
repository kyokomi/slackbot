package slackbot

import (
	"reflect"
	"testing"

	"github.com/alicebob/miniredis"
	redis "gopkg.in/redis.v3"
)

func Test_inMemoryRepository_LoadList(t *testing.T) {
	type fields struct {
		memory map[string]interface{}
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "exist",
			fields: fields{
				memory: map[string]interface{}{
					"load_list": []string{
						"foo",
						"bar",
					},
				},
			},
			args: args{
				key: "load_list",
			},
			want: []string{
				"foo",
				"bar",
			},
			wantErr: false,
		},
		{
			name: "not exist",
			fields: fields{
				memory: map[string]interface{}{},
			},
			args: args{
				key: "empty_load_list",
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &inMemoryRepository{
				memory: tt.fields.memory,
			}
			got, err := r.LoadList(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemoryRepository.LoadList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inMemoryRepository.LoadList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inMemoryRepository_Load(t *testing.T) {
	type fields struct {
		memory map[string]interface{}
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "exist",
			fields: fields{
				memory: map[string]interface{}{
					"load": "foobar",
				},
			},
			args: args{
				key: "load",
			},
			want:    "foobar",
			wantErr: false,
		},
		{
			name: "not exist",
			fields: fields{
				memory: map[string]interface{}{},
			},
			args: args{
				key: "empty_load",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &inMemoryRepository{
				memory: tt.fields.memory,
			}
			got, err := r.Load(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemoryRepository.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("inMemoryRepository.Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisRepository_LoadList(t *testing.T) {
	rc := mockRedis(t)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "exist",
			args: args{
				key: "load_list",
			},
			want: []string{
				"foo",
				"bar",
			},
			wantErr: false,
		},
		{
			name: "not exist",
			args: args{
				key: "empty_load_list",
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.want {
				rc.RPush(tt.args.key, v)
			}
			r := &redisRepository{
				redisDB: rc,
			}
			got, err := r.LoadList(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisRepository.LoadList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisRepository.LoadList() = %v, want %v", got, tt.want)
			}
		})
	}
	rc.Close()
}

func Test_redisRepository_Load(t *testing.T) {
	rc := mockRedis(t)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "exist",
			args: args{
				key: "load",
			},
			want:    "foobar",
			wantErr: false,
		},
		{
			name: "not exist",
			args: args{
				key: "empty_load",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != "" {
				rc.Set(tt.args.key, tt.want, 0)
			}
			r := &redisRepository{
				redisDB: rc,
			}
			got, err := r.Load(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisRepository.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisRepository.Load() = %v, want %v", got, tt.want)
			}
		})
	}
	rc.Close()
}

func mockRedis(t *testing.T) *redis.Client {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("unexpected error while createing test redis server '%#v'", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})
	return client
}
