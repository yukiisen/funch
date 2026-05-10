package cmd

import (
	"errors"
	"github.com/yukiisen/funch/config"
	"reflect"
	"strconv"
	"github.com/spf13/cobra"
)

var SetCmd = &cobra.Command {
	Use: "set <key> <value...>",
	Short: "Modify app settings",
	Long:  "Updates configuration. Supports multi-value input for fields like arguments.",
	Args: cobra.MinimumNArgs(2),
	ValidArgsFunction: getOptions,
	RunE: set,
}

func getOptions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	t := reflect.TypeFor[config.Config]()
	keys := make([]string, 0, t.NumField())

	for k := range t.Fields() {
		keys = append(keys, k.Name)
	}
	
	return keys, cobra.ShellCompDirectiveNoFileComp
}

func set(cmd *cobra.Command, args []string) error {
	key, value := args[0], args[1]
	lib := config.LoadLibrary()

	v := reflect.ValueOf(&lib.Config).Elem()

	field := v.FieldByName(key)

	if !field.IsValid() {
		return errors.New("Unknown Field")
	}

	if !field.CanSet() {
		return errors.New("cannot set field")
	}

	switch field.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)

	case reflect.String:
		field.SetString(value)

	case reflect.Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))

	case reflect.Slice:
		elemType := field.Type().Elem()

		// only support []string for now
		if elemType.Kind() != reflect.String {
			return errors.New("unsupported slice type")
		}

		slice := reflect.MakeSlice(field.Type(), len(args[1:]), len(args[1:]))

		for i, s := range args[1:] {
			slice.Index(i).SetString(s)
		}

		field.Set(slice)



	default:
		return errors.New("unsupported field type")
	}


	err := config.UpdateLibrary(lib);
	if err != nil { return err }

	return nil
}
