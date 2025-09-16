import {useEffect, useState} from "react";

export const ThemeSwitcher = () => {
    const themes = ["emerald", "forest"]
    const [theme, setTheme] = useState(currentTheme())

    function currentTheme(): string {
        return document.documentElement.getAttribute("data-theme") ?? "emerald"
    }

    useEffect(() => {
        const saved = localStorage.getItem("theme");
        if (saved) {
            setTheme(saved);
            document.documentElement.setAttribute("data-theme", saved);
        }
    }, []);

    const toggleTheme = () => {
        const nextIndex = (themes.indexOf(theme) + 1) % themes.length;
        const nextTheme = themes[nextIndex];
        setTheme(nextTheme);
        document.documentElement.setAttribute("data-theme", nextTheme);
        localStorage.setItem("theme", nextTheme);
    }

    return <div className="dropdown dropdown-end">
        <label className="flex cursor-pointer gap-2">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round">
                <circle cx="12" cy="12" r="5"/>
                <path
                    d="M12 1v2M12 21v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M1 12h2M21 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4"/>
            </svg>
            <input
                type="checkbox"
                checked={theme == "forest"}
                className="toggle theme-controller"
                onChange={toggleTheme}
            />
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round">
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
            </svg>
        </label>
    </div>
}
