/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	plists "howett.net/plist"
	"reflect"
	"strings"
)

type Shortcut struct {
	WFWorkflowIcon                      ShortcutIcon
	WFWorkflowActions                   []ShortcutAction
	WFQuickActionSurfaces               []string
	WFWorkflowInputContentItemClasses   []string
	WFWorkflowClientVersion             string
	WFWorkflowMinimumClientVersion      int
	WFWorkflowImportQuestions           interface{}
	WFWorkflowTypes                     []string
	WFWorkflowOutputContentItemClasses  []string
	WFWorkflowHasShortcutInputVariables bool
	WFWorkflowHasOutputFallback         bool
}

type ShortcutIcon struct {
	WFWorkflowIconGlyphNumber int64
	WFWorkflowIconStartColor  int
}

type ShortcutAction struct {
	WFWorkflowActionIdentifier string
	WFWorkflowActionParameters map[string]any
}

var data Shortcut
var code strings.Builder

func decompile(b []byte) {
	var _, marshalErr = plists.Unmarshal(b, &data)
	handle(marshalErr)

	decompileIcon()
	decompileActions()
}

func decompileIcon() {
	var icon = data.WFWorkflowIcon
	if icon.WFWorkflowIconStartColor != iconColor {
		makeColors()
		for name, i := range colors {
			if icon.WFWorkflowIconStartColor != i {
				continue
			}

			code.WriteString(fmt.Sprintf("#define color %s", name))
		}
	}

	if icon.WFWorkflowIconGlyphNumber != iconGlyph {
		makeGlyphs()
		for name, i := range glyphs {
			if icon.WFWorkflowIconGlyphNumber != i {
				continue
			}

			code.WriteString(fmt.Sprintf("#define glyph %s\n", name))
		}
	}

	code.WriteRune('\n')
}

func decompileActions() {
	var variableValue any
	for _, action := range data.WFWorkflowActions {
		var identifier = matchAction(action.WFWorkflowActionIdentifier)
		if identifier == "" {
			switch action.WFWorkflowActionIdentifier {
			case "is.workflow.actions.gettext":
				var text = action.WFWorkflowActionParameters["WFTextActionText"].(map[string]any)

				if reflect.TypeOf(text).String() == "string" {
					variableValue = fmt.Sprintf("\"%v\"", text)
				}
			case "is.workflow.actions.number":
				variableValue = action.WFWorkflowActionParameters["WFNumberActionNumber"].(string)
			case "is.workflow.actions.setvariable":
				code.WriteRune('@')
				code.WriteString(action.WFWorkflowActionParameters["WFVariableName"].(string))

				if variableValue != nil {
					code.WriteRune(' ')
					code.WriteRune('=')
					code.WriteRune(' ')
					code.WriteString(fmt.Sprintf("%v", variableValue))
				}

				variableValue = nil
				code.WriteRune('\n')
			}
			continue
		}

		code.WriteString(identifier)
		code.WriteRune('(')

		code.WriteString(")\n")
	}
}

func matchAction(identifier string) string {
	standardActions()
	for action, data := range actions {
		if data.identifier == identifier {
			return action
		}
	}

	return ""
}
