module github.com/twiglab/h2o/chrgg

go 1.26.0

require (
	entgo.io/ent v0.14.5
	github.com/duckdb/duckdb-go/v2 v2.5.5
	github.com/eclipse/paho.mqtt.golang v1.5.1
	github.com/go-sql-driver/mysql v1.9.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/mattn/go-sqlite3 v1.14.34
	github.com/spf13/cobra v1.10.2
	github.com/spf13/viper v1.21.0
	github.com/twiglab/h2o v0.0.0-00010101000000-000000000000
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	ariga.io/atlas v0.32.1-0.20250325101103-175b25e1c1b9 // indirect
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/99designs/gqlgen v0.17.86 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/apache/arrow-go/v18 v18.5.1 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/duckdb/duckdb-go-bindings v0.3.3 // indirect
	github.com/duckdb/duckdb-go-bindings/lib/darwin-amd64 v0.3.3 // indirect
	github.com/duckdb/duckdb-go-bindings/lib/darwin-arm64 v0.3.3 // indirect
	github.com/duckdb/duckdb-go-bindings/lib/linux-amd64 v0.3.3 // indirect
	github.com/duckdb/duckdb-go-bindings/lib/linux-arm64 v0.3.3 // indirect
	github.com/duckdb/duckdb-go-bindings/lib/windows-amd64 v0.3.3 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/google/flatbuffers v25.12.19+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/hcl/v2 v2.18.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.25 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/urfave/cli/v3 v3.6.1 // indirect
	github.com/vektah/gqlparser/v2 v2.5.31 // indirect
	github.com/zclconf/go-cty v1.14.4 // indirect
	github.com/zclconf/go-cty-yaml v1.1.0 // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/exp v0.0.0-20260112195511-716be5621a96 // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/telemetry v0.0.0-20260209163413-e7419c687ee4 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/tools v0.42.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
)

tool (
	entgo.io/ent/cmd/ent
	github.com/99designs/gqlgen
)

replace github.com/twiglab/h2o => ../
