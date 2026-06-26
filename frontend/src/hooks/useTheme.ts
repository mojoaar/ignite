import { useEffect } from "react";
import { useThemeStore } from "@/lib/store/theme";

export function useTheme() {
  const mode = useThemeStore((s) => s.mode);
  const toggle = useThemeStore((s) => s.toggle);
  const init = useThemeStore((s) => s.init);
  useEffect(() => { init(); }, [init]);
  return { mode, toggle };
}
