import {createContext} from "react";

export type Theme = "corporate" | "business";

export interface ThemeContextType {
    theme: Theme;
    switchTheme: () => void;
}

export const ThemeContext = createContext<ThemeContextType | undefined>(undefined);
