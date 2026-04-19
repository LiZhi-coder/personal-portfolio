import Footer from "@/components/Footer";
import Navigation from "@/components/Navigation";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useFetchData } from "@/hooks/useFetchData";
import { Helmet } from "react-helmet-async";
import { useState } from "react";

interface PhotoItem {
  id: string;
  title: string;
  date: string;
  description: string;
  image: string;
  alt: string;
}

interface PhotoVisualProps {
  photo: PhotoItem;
  className?: string;
  fallbackClassName?: string;
}

function formatDate(date: string) {
  const parsed = new Date(date);

  if (Number.isNaN(parsed.getTime())) {
    return date;
  }

  return parsed.toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

function PhotoVisual({
  photo,
  className = "h-full w-full object-cover",
  fallbackClassName = "flex h-full w-full items-center justify-center bg-muted/70 px-6 text-center text-sm text-muted-foreground",
}: PhotoVisualProps) {
  const [hasError, setHasError] = useState(false);

  if (hasError) {
    return (
      <div className={fallbackClassName} role="img" aria-label={photo.alt}>
        <span>{photo.alt}</span>
      </div>
    );
  }

  return (
    <img
      src={photo.image}
      alt={photo.alt}
      className={className}
      loading="lazy"
      onError={() => setHasError(true)}
    />
  );
}

export default function Photos() {
  const {
    data: photos,
    loading,
    error,
  } = useFetchData<PhotoItem>("/photos.json", "加载照片列表失败");
  const [selectedPhoto, setSelectedPhoto] = useState<PhotoItem | null>(null);

  return (
    <div className="min-h-screen flex flex-col bg-background font-sans selection:bg-primary/10">
      <Helmet>
        <title>照片墙 - 生活切片</title>
        <meta
          name="description"
          content="收集一些路过的风景、光影和生活片段。"
        />
      </Helmet>

      <Navigation />

      <main className="flex-1">
        <section className="py-20 px-6 sm:px-8">
          <div className="max-w-6xl mx-auto space-y-12">
            <div
              className="space-y-4 animate-fade-in-up"
              style={{ animationDelay: "0.1s" }}
            >
              <h1 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight">
                照片墙
              </h1>
              <p className="text-lg text-muted-foreground font-light max-w-2xl">
                一些瞬间，来自通勤、散步、路上和抬头时看到的光。
              </p>
            </div>

            {loading && (
              <div className="flex items-center justify-center py-20">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
              </div>
            )}

            {error && (
              <div className="text-center py-20">
                <p className="text-destructive text-lg">{error}</p>
              </div>
            )}

            {!loading && !error && photos.length === 0 && (
              <div className="text-center py-20">
                <p className="text-muted-foreground text-lg">暂无照片</p>
              </div>
            )}

            {!loading && !error && photos.length > 0 && (
              <div
                className="grid grid-cols-1 gap-6 sm:grid-cols-2 xl:grid-cols-3 animate-fade-in-up"
                style={{ animationDelay: "0.2s" }}
              >
                {photos.map((photo) => (
                  <button
                    key={photo.id}
                    type="button"
                    onClick={() => setSelectedPhoto(photo)}
                    className="group overflow-hidden rounded-2xl border border-border/40 bg-card text-left transition-all duration-300 hover:-translate-y-1 hover:border-border/60 hover:bg-card/80 hover:shadow-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  >
                    <div className="aspect-[4/3] overflow-hidden bg-muted">
                      <PhotoVisual
                        photo={photo}
                        className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-[1.03]"
                      />
                    </div>

                    <div className="space-y-3 p-5">
                      <div className="flex items-start justify-between gap-4">
                        <h3 className="text-lg font-semibold tracking-tight text-foreground group-hover:text-primary transition-colors">
                          {photo.title}
                        </h3>
                        <time className="shrink-0 whitespace-nowrap rounded-full border border-border/20 bg-muted/50 px-3 py-1 text-xs text-muted-foreground">
                          {formatDate(photo.date)}
                        </time>
                      </div>

                      <p className="line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                        {photo.description}
                      </p>
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>
        </section>
      </main>

      <Dialog
        open={selectedPhoto !== null}
        onOpenChange={(open) => {
          if (!open) {
            setSelectedPhoto(null);
          }
        }}
      >
        {selectedPhoto && (
          <DialogContent className="max-h-[90vh] overflow-y-auto border-border/60 bg-background p-4 sm:max-w-4xl sm:p-6">
            <div className="space-y-5">
              <div className="overflow-hidden rounded-xl border border-border/40 bg-muted">
                <PhotoVisual
                  key={`modal-${selectedPhoto.id}`}
                  photo={selectedPhoto}
                  className="max-h-[70vh] w-full object-contain bg-card"
                  fallbackClassName="flex min-h-72 w-full items-center justify-center bg-muted/70 px-8 py-12 text-center text-base text-muted-foreground"
                />
              </div>

              <DialogHeader className="space-y-3 text-left">
                <div className="flex flex-wrap items-center gap-3">
                  <DialogTitle className="text-2xl font-semibold tracking-tight text-foreground">
                    {selectedPhoto.title}
                  </DialogTitle>
                  <time className="rounded-full border border-border/20 bg-muted/50 px-3 py-1 text-xs text-muted-foreground">
                    {formatDate(selectedPhoto.date)}
                  </time>
                </div>
                <DialogDescription className="text-sm leading-relaxed text-muted-foreground sm:text-base">
                  {selectedPhoto.description}
                </DialogDescription>
              </DialogHeader>
            </div>
          </DialogContent>
        )}
      </Dialog>

      <Footer />
    </div>
  );
}
