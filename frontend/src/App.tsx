import { useEffect, useState } from "react";
import { useTheme } from "@/hooks/useTheme";
import { Greet } from "@wails/go/main/App";

function App() {
  const [message, setMessage] = useState("");
  const { mode, toggle } = useTheme();
  useEffect(() => { Greet("world").then(setMessage); }, []);
  return (
    <div className="flex h-screen flex-col items-center justify-center gap-4 bg-background text-text-primary">
      <h1 className="font-mono text-2xl">{message}</h1>
      <button onClick={toggle}
        className="rounded-md border border-border bg-surface px-4 py-2 font-mono text-sm text-text-secondary hover:bg-surface-hover">
        {mode === "dark" ? "Light" : "Dark"} mode
      </button>
    </div>
  );
}
export default App;
