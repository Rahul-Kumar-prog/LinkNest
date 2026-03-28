import { useEffect, useMemo, useRef, useState } from "react";

const platformMeta = {
    x: {
        label: "X",
    },
    linkedin: {
        label: "LinkedIn",
        badge: "in",
    },
};

function MediaIcon() {
    return (
        <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <rect x="3.5" y="5" width="17" height="14" rx="2.5" />
            <circle cx="9" cy="10" r="1.5" />
            <path d="M5.5 16l4.5-4.5 3.5 3.5 2.5-2.5 2.5 3" />
        </svg>
    );
}

function ChevronIcon({ open }) {
    return (
        <svg viewBox="0 0 20 20" className={`h-4 w-4 transition ${open ? "rotate-180" : ""}`} fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <path d="M5 7.5 10 12.5l5-5" />
        </svg>
    );
}

function InfoIcon() {
    return (
        <svg viewBox="0 0 20 20" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="10" cy="10" r="7" />
            <path d="M10 8v4" />
            <circle cx="10" cy="5.5" r="0.8" fill="currentColor" stroke="none" />
        </svg>
    );
}

function CommentIcon() {
    return (
        <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <path d="M7 18.5c-1.7 0-3-1.3-3-3v-7c0-1.7 1.3-3 3-3h10c1.7 0 3 1.3 3 3v7c0 1.7-1.3 3-3 3H11l-4 3v-3H7Z" />
        </svg>
    );
}

function RepostIcon() {
    return (
        <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <path d="M7 7h8l-2.5-2.5" />
            <path d="M17 17H9l2.5 2.5" />
            <path d="M15 4.5 18.5 8 15 11.5" />
            <path d="M9 12.5 5.5 16 9 19.5" />
        </svg>
    );
}

function LikeIcon() {
    return (
        <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <path d="M12 20s-6.5-3.9-8.4-7.6C2.2 9.8 3.5 6.5 7 6.5c2.1 0 3.3 1.2 4 2.3.7-1.1 1.9-2.3 4-2.3 3.5 0 4.8 3.3 3.4 5.9C18.5 16.1 12 20 12 20Z" />
        </svg>
    );
}

function ShareIcon() {
    return (
        <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <path d="M14 5h5v5" />
            <path d="M10 14 19 5" />
            <path d="M19 13v5a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h5" />
        </svg>
    );
}

function AnalyticsIcon() {
    return (
        <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <path d="M5 18V9" />
            <path d="M12 18V5" />
            <path d="M19 18v-7" />
        </svg>
    );
}

function BookmarkIcon() {
    return (
        <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <path d="M7 5.5h10a1 1 0 0 1 1 1v13l-6-3-6 3v-13a1 1 0 0 1 1-1Z" />
        </svg>
    );
}

function XLogo({ className = "h-4 w-4" }) {
    return (
        <svg viewBox="0 0 24 24" className={className} fill="currentColor" aria-hidden="true">
            <path d="M18.901 2H22l-6.768 7.736L23.194 22H16.96l-4.884-7.443L5.56 22H2.458l7.237-8.272L1.806 2H8.2l4.415 6.73L18.901 2Zm-1.087 18.112h1.717L7.267 3.79H5.425l12.389 16.322Z" />
        </svg>
    );
}

function PlatformToggleIcon({ platform, selected, connected, onClick }) {
    const meta = platformMeta[platform];
    const baseClass = selected
        ? "border-sky-300 bg-sky-400/15 text-white"
        : "border-white/10 bg-white/[0.04] text-stone-300 hover:bg-white/[0.08]";

    return (
        <button
            type="button"
            onClick={onClick}
            className={`flex h-14 w-14 items-center justify-center rounded-full border text-sm font-bold transition ${baseClass} ${!connected ? "opacity-60" : ""}`}
        >
            {platform === "x" ? <XLogo className="h-4 w-4" /> : meta.badge}
        </button>
    );
}

function PreviewAvatar({ image, fallback, className = "h-12 w-12" }) {
    if (image) {
        return <img src={image} alt={fallback} className={`${className} rounded-full object-cover`} />;
    }

    return (
        <div className={`${className} flex items-center justify-center rounded-full bg-slate-700 text-sm font-bold text-white`}>
            {fallback.slice(0, 2).toUpperCase()}
        </div>
    );
}

function LinkedInPreview({ draft, attachedMedia, user }) {
    const [expanded, setExpanded] = useState(false);
    const linkedInLimit = 280;
    const hasDraft = Boolean(draft);
    const isTruncated = hasDraft && draft.length > linkedInLimit;
    const previewText = hasDraft
        ? (expanded ? draft : draft.slice(0, linkedInLimit))
        : "Your LinkedIn post preview will appear here.";

    useEffect(() => {
        if (!isTruncated) {
            setExpanded(false);
        }
    }, [isTruncated, draft]);

    return (
        <section className="min-w-0 rounded-[2rem] border border-white/10 bg-slate-900/80 p-6">
            <div className="flex items-center gap-2 text-sm font-semibold text-white">
                <span className="flex h-8 w-8 items-center justify-center rounded-full bg-[#0A66C2] text-sm font-bold text-white">in</span>
                <span className="text-stone-400"><InfoIcon /></span>
            </div>
            <div className="mt-4 w-full min-w-0 overflow-hidden rounded-[1.2rem] border border-stone-200 bg-white p-5 text-slate-900 shadow-lg shadow-black/10">
                <div className="flex items-start gap-3">
                    <PreviewAvatar
                        image={user?.picture}
                        fallback={user?.name || user?.email || "LI"}
                        className="h-14 w-14"
                    />
                    <div className="min-w-0">
                        <p className="text-xl font-semibold leading-6 text-slate-900">{user?.name || "LinkedIn User"}</p>
                        <p className="mt-1 text-sm text-slate-500">1h · 🌐</p>
                    </div>
                </div>
                <div className="mt-4">
                    <p className="whitespace-pre-wrap break-words [overflow-wrap:anywhere] text-base leading-7 text-slate-900">
                        {previewText}
                    </p>
                    {isTruncated ? (
                        <div className="mt-1 text-right">
                            <button
                                type="button"
                                onClick={() => setExpanded((current) => !current)}
                                className="cursor-pointer text-sm font-medium text-slate-500 transition hover:text-slate-700"
                            >
                                {expanded ? "show less" : "...more"}
                            </button>
                        </div>
                    ) : null}
                </div>
                {attachedMedia ? (
                    <img
                        src={attachedMedia.previewUrl}
                        alt={attachedMedia.name}
                        className="mt-4 max-h-72 w-full rounded-[1rem] object-cover"
                    />
                ) : null}
                <div className="mt-4 border-t border-slate-200 pt-4 text-sm text-slate-700">
                    <div className="grid grid-cols-4 gap-3 text-center font-semibold">
                        <div className="flex items-center justify-center gap-2">
                            <LikeIcon />
                            <span>Like</span>
                        </div>
                        <div className="flex items-center justify-center gap-2">
                            <CommentIcon />
                            <span>Comment</span>
                        </div>
                        <div className="flex items-center justify-center gap-2">
                            <RepostIcon />
                            <span>Repost</span>
                        </div>
                        <div className="flex items-center justify-center gap-2">
                            <ShareIcon />
                            <span>Send</span>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
}


function HighlightedComposer({ draft, setDraft }) {
    return (
        <textarea
            value={draft}
            onChange={(event) => setDraft(event.target.value)}
            placeholder="Share a launch update, a product insight, or your next announcement."
            className="h-80 w-full rounded-[1.5rem] border border-white/10 bg-white/[0.03] p-6 text-lg text-white outline-none transition placeholder:text-stone-500 focus:border-sky-400/60 break-words"
        />
    );
}

function XPreview({ draft, attachedMedia, user, username, xLimit }) {
    const previewText = draft ? draft.slice(0, xLimit) : "Your X post preview will appear here.";

    return (
        <section className="min-w-0 rounded-[2rem] border border-white/10 bg-slate-900/80 p-6">
            <div className="flex items-center gap-2 text-sm font-semibold text-white">
                <span className="flex h-8 w-8 items-center justify-center rounded-full bg-white text-slate-950"><XLogo className="h-4 w-4" /></span>
                <span className="text-stone-400"><InfoIcon /></span>
            </div>
            <div className="mt-4 w-full min-w-0 overflow-hidden rounded-[1.2rem] border border-stone-200 bg-white p-5 text-slate-900 shadow-lg shadow-black/10">
                <div className="flex items-start gap-3">
                    <PreviewAvatar
                        image={user?.picture}
                        fallback={user?.name || user?.email || "X"}
                        className="h-12 w-12"
                    />
                    <div className="min-w-0">
                        <p className="text-lg font-semibold leading-5 text-slate-900">
                            {user?.name || "X User"} <span className="font-normal text-slate-600">@{username || "username"}</span>
                        </p>
                    </div>
                </div>
                <p className="mt-4 whitespace-pre-wrap break-words [overflow-wrap:anywhere] text-base leading-7 text-slate-900">
                    {previewText}
                </p>
                {attachedMedia ? (
                    <img
                        src={attachedMedia.previewUrl}
                        alt={attachedMedia.name}
                        className="mt-4 max-h-72 w-full rounded-[1rem] border border-stone-200 object-cover"
                    />
                ) : null}
                <div className="mt-4 border-t border-slate-200 pt-4 text-sm text-slate-600">
                    <div className="grid grid-cols-6 gap-3 text-center">
                        <div className="flex items-center justify-center"><CommentIcon /></div>
                        <div className="flex items-center justify-center"><RepostIcon /></div>
                        <div className="flex items-center justify-center"><LikeIcon /></div>
                        <div className="flex items-center justify-center"><AnalyticsIcon /></div>
                        <div className="flex items-center justify-center"><BookmarkIcon /></div>
                        <div className="flex items-center justify-center"><ShareIcon /></div>
                    </div>
                </div>
            </div>
        </section>
    );
}

export default function Hero({
    draft,
    setDraft,
    attachedMedia,
    selectedPlatforms,
    togglePlatform,
    platforms,
    notice,
    publishResults,
    onPublish,
    onPickMedia,
    onRemoveMedia,
    publishing,
    xLimit,
    user,
}) {
    const [pickerOpen, setPickerOpen] = useState(false);
    const pickerRef = useRef(null);
    const xCount = draft.length;
    const xTooLong = xCount > xLimit;
    const mediaPreviewUrl = useMemo(() => {
        if (!attachedMedia) {
            return "";
        }

        return URL.createObjectURL(attachedMedia);
    }, [attachedMedia]);

    useEffect(() => {
        return () => {
            if (mediaPreviewUrl) {
                URL.revokeObjectURL(mediaPreviewUrl);
            }
        };
    }, [mediaPreviewUrl]);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (pickerRef.current && !pickerRef.current.contains(event.target)) {
                setPickerOpen(false);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const previewMedia = attachedMedia
        ? {
            ...attachedMedia,
            previewUrl: mediaPreviewUrl,
        }
        : null;

    const showLinkedInPreview = selectedPlatforms.includes("linkedin");
    const showXPreview = selectedPlatforms.includes("x");

    return (
        <div className="min-h-screen bg-[linear-gradient(180deg,_rgba(12,10,9,0.96)_0%,_rgba(17,24,39,0.92)_50%,_rgba(2,6,23,1)_100%)] px-6 py-8 lg:px-10">
            <div className="mx-auto max-w-6xl space-y-8">
                <header className="rounded-[2rem] border border-white/10 bg-white/5 p-8 shadow-2xl shadow-black/20 backdrop-blur">
                    <p className="text-sm uppercase tracking-[0.3em] text-stone-400">Compose</p>
                    <div className="mt-4 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
                        <div className="max-w-2xl">
                            <h1 className="text-4xl font-black text-white">Publish to your connected channels</h1>
                            <p className="mt-3 text-base text-stone-300">
                                Choose X, LinkedIn, or both. This MVP supports text posts, plus one attached image in All Channels.
                            </p>
                        </div>
                        <div className={`rounded-full border px-4 py-2 text-sm ${xTooLong ? "border-amber-400/40 bg-amber-400/10 text-amber-100" : "border-emerald-400/30 bg-emerald-400/10 text-emerald-100"}`}>
                            X count {xCount}/{xLimit}
                        </div>
                    </div>
                </header>

                <div className="grid items-start gap-8 xl:grid-cols-[minmax(0,1.35fr)_minmax(0,0.85fr)]">
                    <section className="min-w-0 self-start rounded-[2rem] border border-white/10 bg-slate-950/80 p-6 shadow-2xl shadow-black/30">
                        <div className="mb-4 flex items-center justify-between gap-4">
                            <h2 className="text-xl font-semibold text-white">Post text</h2>
                            <div className="flex items-center gap-3" ref={pickerRef}>
                                <div className="relative">
                                    <button
                                        type="button"
                                        onClick={() => setPickerOpen((current) => !current)}
                                        className="inline-flex items-center gap-3 rounded-full border border-white/10 px-4 py-2 text-xs font-semibold uppercase tracking-[0.2em] text-stone-300 transition hover:bg-white/5"
                                    >
                                        <span>Channels</span>
                                        <ChevronIcon open={pickerOpen} />
                                    </button>
                                    {pickerOpen ? (
                                        <div className="absolute right-0 top-[calc(100%+0.75rem)] z-20 rounded-[1.5rem] border border-white/10 bg-slate-950/95 p-4 shadow-2xl shadow-black/50 backdrop-blur">
                                            <div className="flex items-center gap-3">
                                                {Object.keys(platformMeta).map((platform) => {
                                                    const connected = !!platforms?.[platform]?.connected && !(platform === "x" && xTooLong);
                                                    return (
                                                        <PlatformToggleIcon
                                                            key={platform}
                                                            platform={platform}
                                                            selected={selectedPlatforms.includes(platform)}
                                                            connected={connected}
                                                            onClick={() => togglePlatform(platform)}
                                                        />
                                                    );
                                                })}
                                            </div>
                                        </div>
                                    ) : null}
                                </div>
                                <button
                                    type="button"
                                    onClick={() => {
                                        setDraft("");
                                        onRemoveMedia();
                                    }}
                                    className="rounded-full border border-white/10 px-4 py-2 text-xs font-semibold uppercase tracking-[0.2em] text-stone-300 transition hover:bg-white/5"
                                >
                                    Clear
                                </button>
                            </div>
                        </div>
                        <HighlightedComposer
                            draft={draft}
                            setDraft={setDraft}
                        />
                        <div className="mt-4 flex items-end justify-between gap-4 border-t border-white/10 pt-4">
                            <div className="flex items-center gap-3">
                                <label className="inline-flex cursor-pointer items-center justify-center rounded-xl border border-sky-400/30 bg-sky-400/10 p-3 text-sky-200 transition hover:bg-sky-400/20">
                                    <MediaIcon />
                                    <input
                                        type="file"
                                        accept="image/png,image/jpeg,image/jpg,image/gif,image/webp"
                                        className="hidden"
                                        onChange={(event) => onPickMedia(event.target.files?.[0] || null)}
                                    />
                                </label>
                                <p className="text-sm text-stone-500">
                                    {attachedMedia ? attachedMedia.name : "Attach one image"}
                                </p>
                            </div>
                            <button
                                type="button"
                                onClick={onPublish}
                                disabled={publishing}
                                className="rounded-full bg-sky-400 px-8 py-3 text-sm font-semibold text-slate-950 transition hover:bg-sky-300 disabled:cursor-not-allowed disabled:bg-sky-400/60"
                            >
                                {publishing ? "Publishing..." : "Publish now"}
                            </button>
                        </div>
                        {attachedMedia ? (
                            <div className="mt-5 rounded-[1.5rem] border border-white/10 bg-white/[0.04] p-4">
                                <div className="flex items-start justify-between gap-4">
                                    <div>
                                        <p className="text-sm font-semibold text-white">Attached image</p>
                                        <p className="mt-1 text-xs text-stone-400">{attachedMedia.name}</p>
                                    </div>
                                    <button
                                        type="button"
                                        onClick={onRemoveMedia}
                                        className="rounded-full border border-rose-400/30 bg-rose-400/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-rose-100 transition hover:bg-rose-400/20"
                                    >
                                        Remove
                                    </button>
                                </div>
                                <img
                                    src={previewMedia?.previewUrl}
                                    alt={attachedMedia.name}
                                    className="mt-4 max-h-72 w-full rounded-[1.25rem] object-cover"
                                />
                            </div>
                        ) : null}
                    </section>

                    <aside className="min-w-0 space-y-6">
                        {showLinkedInPreview ? (
                            <LinkedInPreview
                                draft={draft}
                                attachedMedia={previewMedia}
                                user={{
                                    name: platforms?.linkedin?.name || user?.name,
                                    picture: platforms?.linkedin?.profile_picture || user?.picture,
                                    email: user?.email,
                                }}
                            />
                        ) : null}
                        {showXPreview ? (
                            <XPreview
                                draft={draft}
                                attachedMedia={previewMedia}
                                user={{
                                    name: platforms?.x?.name || user?.name,
                                    picture: platforms?.x?.profile_picture || user?.picture,
                                    email: user?.email,
                                }}
                                username={platforms?.x?.username}
                                xLimit={xLimit}
                            />
                        ) : null}
                        <section className="rounded-[2rem] border border-white/10 bg-slate-900/80 p-6">
                            <h2 className="text-xl font-semibold text-white">Status</h2>
                            {notice ? (
                                <div className="mt-4 rounded-2xl border border-white/10 bg-white/[0.04] p-4 text-sm text-stone-200">
                                    {notice}
                                </div>
                            ) : (
                                <p className="mt-4 text-sm text-stone-400">Connection updates and publish results will appear here.</p>
                            )}

                            {publishResults.length > 0 && (
                                <div className="mt-4 space-y-3">
                                    {publishResults.map((result) => (
                                        <div
                                            key={result.platform}
                                            className={`rounded-2xl border p-4 text-sm ${
                                                result.success
                                                    ? "border-emerald-400/30 bg-emerald-400/10 text-emerald-100"
                                                    : "border-rose-400/30 bg-rose-400/10 text-rose-100"
                                            }`}
                                        >
                                            <span className="font-semibold">{platformMeta[result.platform]?.label || result.platform}:</span>{" "}
                                            {result.success ? "Published successfully." : result.error}
                                        </div>
                                    ))}
                                </div>
                            )}
                        </section>
                    </aside>
                </div>
            </div>
        </div>
    );
}
