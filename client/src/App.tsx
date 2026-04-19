import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import BlogDetail from "@/pages/BlogDetail";
import NotFound from "@/pages/NotFound";
import Photos from "@/pages/Photos";
import ProjectDetail from "@/pages/ProjectDetail";
import { HelmetProvider } from "react-helmet-async";
import { Route, Switch, useLocation } from "wouter";
import { AnimatePresence, motion } from "framer-motion";
import ErrorBoundary from "./components/ErrorBoundary";
import { ThemeProvider } from "./contexts/ThemeContext";
import Home from "./pages/Home";
import Projects from "./pages/Projects";
import Resume from "./pages/Resume";
import Blog from "./pages/Blog";


function Router() {
  const [location] = useLocation();

  return (
    <AnimatePresence>
      <motion.div
        key={location}
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: -10 }}
        transition={{ duration: 0.15 }}
      >
        <Switch>
          <Route path="/" component={Home} />
          <Route path="/projects" component={Projects} />
          <Route path="/projects/:id" component={ProjectDetail} />
          <Route path="/photos" component={Photos} />
          <Route path="/resume" component={Resume} />
          <Route path="/blog" component={Blog} />
          <Route path="/blog/:id" component={BlogDetail} />
          <Route path="/404" component={NotFound} />
          {/* Final fallback route */}
          <Route component={NotFound} />
        </Switch>
      </motion.div>
    </AnimatePresence>
  );
}

// NOTE: About Theme
// - First choose a default theme according to your design style (dark or light bg), than change color palette in index.css
//   to keep consistent foreground/background color across components
// - If you want to make theme switchable, pass `switchable` ThemeProvider and use `useTheme` hook

function App() {
  return (
    <ErrorBoundary>
      <HelmetProvider>
        <ThemeProvider
          defaultTheme="light"
        >
          <TooltipProvider>
            <Toaster />
            <Router />
          </TooltipProvider>
        </ThemeProvider>
      </HelmetProvider>
    </ErrorBoundary>
  );
}

export default App;
