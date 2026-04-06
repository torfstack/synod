import {useContext} from "react";
import {ThemeContext, type ThemeContextType} from "./theme-context.ts";

export const useTheme = (): ThemeContextType => {
    const ctx = useContext(ThemeContext);
    if (!ctx) throw new Error("useTheme must be used within ThemeProvider");
    return ctx;
};
