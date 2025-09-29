import React, {createContext, useContext, useEffect, useState} from "react";

type Theme = "corporate" | "business"

interface ThemeContextType {
    theme: Theme;
    switchTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const ThemeProvider = ({children}: { children: React.ReactNode }) => {
    const getPreferredTheme = (): Theme => {
        const stored = localStorage.getItem("theme") as Theme | null;
        if (stored) return stored;

        if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
            return "business";
        }

        return "corporate";
    };

    const [theme, setTheme] = useState<Theme>(getPreferredTheme);

    useEffect(() => {
        document.documentElement.setAttribute("data-theme", theme);
        localStorage.setItem("theme", theme);
    }, [theme]);

    useEffect(() => {
        const media = window.matchMedia("(prefers-color-scheme: dark)");
        const handler = (e: MediaQueryListEvent) => {
            if (!localStorage.getItem("theme")) {
                setTheme(e.matches ? "business" : "corporate");
            }
        };
        media.addEventListener("change", handler);
        return () => media.removeEventListener("change", handler);
    }, []);

    const switchTheme = () => {
        setTheme(currentTheme => {
            switch (currentTheme) {
                case "corporate":
                    return "business"
                case "business":
                    return "corporate"
            }
        })
    }

    return (
        <ThemeContext.Provider value={{theme, switchTheme}}>
            {children}
        </ThemeContext.Provider>
    );
};

export const useTheme = () => {
    const ctx = useContext(ThemeContext);
    if (!ctx) throw new Error("useTheme must be used within ThemeProvider");
    return ctx;
};
