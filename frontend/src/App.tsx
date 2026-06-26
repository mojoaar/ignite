import { useState, useEffect } from "react";
import { Greet } from "@wails/go/main/App";

function App() {
  const [message, setMessage] = useState("");
  useEffect(() => { Greet("world").then(setMessage); }, []);
  return (
    <div className="flex h-screen items-center justify-center bg-background text-text-primary">
      <h1 className="font-mono text-2xl">{message}</h1>
    </div>
  );
}
export default App;
