import { Component, type ReactNode } from "react";

interface Props {
  children: ReactNode;
}

interface State {
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  render() {
    if (this.state.error) {
      return (
        <div className="flex h-screen flex-col items-center justify-center gap-4 bg-background text-text-primary p-8">
          <h1 className="font-mono text-xl font-semibold">Something went wrong</h1>
          <p className="font-mono text-sm text-text-secondary text-center max-w-md">
            {this.state.error.message || "An unexpected error occurred."}
          </p>
          <button
            onClick={() => {
              this.setState({ error: null });
              window.location.reload();
            }}
            className="rounded-lg bg-accent px-4 py-2 font-mono text-sm text-white hover:bg-accent-light transition-colors"
          >
            Reload
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}
