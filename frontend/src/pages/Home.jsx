import { useEffect, useMemo, useState } from "react";
import Hero from "../Components/Landing/Hero";
import Sidebar from "../Components/Sidebar/Sidebar";

const API_BASE = "http://localhost:8080";
const platformRoutes = {
    x: "twitter",
    linkedin: "linkedin",
};

export default function Home() {
    const [session, setSession] = useState(null);
    const [loading, setLoading] = useState(true);
    const [draft, setDraft] = useState("");
    const [attachedMedia, setAttachedMedia] = useState(null);
    const [selectedPlatforms, setSelectedPlatforms] = useState([]);
    const [publishResults, setPublishResults] = useState([]);
    const [notice, setNotice] = useState("");
    const [publishing, setPublishing] = useState(false);
    const [busyPlatform, setBusyPlatform] = useState("");
    const [activeChannel, setActiveChannel] = useState("all");

    const xLimit = session?.limits?.x_post_length ?? 280;
    const xTooLong = draft.length > xLimit;

    const refreshStatus = async () => {
        try {
            const response = await fetch(`${API_BASE}/api/me/status`, {
                credentials: "include",
            });

            if (response.status === 401) {
                setSession(null);
                return;
            }

            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.error || "Failed to load account status");
            }

            setSession(data);
        } catch (error) {
            console.error(error);
            setNotice(error.message);
            setSession(null);
        }
    };

    useEffect(() => {
        const syncSession = async () => {
            const params = new URLSearchParams(window.location.search);
            const auth = params.get("auth");
            const oauth = params.get("oauth");
            const provider = params.get("provider");
            const platform = params.get("platform");

            if (auth === "success" && provider === "google") {
                setNotice("Signed in with Google.");
            } else if (auth === "error" && provider === "google") {
                setNotice("Google sign in failed.");
            } else if (oauth === "success" && platform) {
                setNotice(`${platform === "x" ? "X" : "LinkedIn"} connected successfully.`);
            } else if (oauth === "error" && platform) {
                setNotice(`Failed to connect ${platform === "x" ? "X" : "LinkedIn"}.`);
            }

            if (params.toString()) {
                window.history.replaceState({}, document.title, window.location.pathname);
            }

            await refreshStatus();
            setLoading(false);
        };

        syncSession();
    }, []);

    const connectedPlatforms = useMemo(() => ({
        x: !!session?.platforms?.x?.connected,
        linkedin: !!session?.platforms?.linkedin?.connected,
    }), [session]);

    const handleGoogleLogin = () => {
        window.location.href = `${API_BASE}/api/auth/google`;
    };

    const handlePickMedia = (file) => {
        if (!file) {
            setAttachedMedia(null);
            return;
        }

        if (!file.type.startsWith("image/")) {
            setNotice("Only image attachments are supported right now.");
            return;
        }

        setNotice("");
        setAttachedMedia(file);
    };

    const handleRemoveMedia = () => {
        setAttachedMedia(null);
    };

    const handleConnectPlatform = (platform) => {
        setBusyPlatform(platform);
        window.location.href = `${API_BASE}/api/${platformRoutes[platform] || platform}`;
    };

    const handleDisconnectPlatform = async (platform) => {
        setBusyPlatform(platform);
        try {
            const response = await fetch(`${API_BASE}/api/platforms/disconnect?platform=${platform}`, {
                method: "POST",
                credentials: "include",
            });
            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.error || `Failed to disconnect ${platform}`);
            }
            setNotice(data.message || `${platform === "x" ? "X" : "LinkedIn"} disconnected.`);
            setSelectedPlatforms((current) => current.filter((item) => item !== platform));
            await refreshStatus();
        } catch (error) {
            console.error(error);
            setNotice(error.message);
        } finally {
            setBusyPlatform("");
        }
    };

    const handleOpenChannel = (channelId) => {
        setActiveChannel(channelId);
    };

    const handleLogout = async () => {
        try {
            await fetch(`${API_BASE}/api/logout`, {
                method: "POST",
                credentials: "include",
            });
        } finally {
            localStorage.removeItem("token");
            setSession(null);
            setSelectedPlatforms([]);
            setPublishResults([]);
            setAttachedMedia(null);
        }
    };

    const togglePlatform = (platform) => {
        if (!connectedPlatforms[platform]) {
            setNotice(`Connect ${platform === "x" ? "X" : "LinkedIn"} before publishing.`);
            return;
        }

        setSelectedPlatforms((current) =>
            current.includes(platform)
                ? current.filter((item) => item !== platform)
                : [...current, platform]
        );
    };

    const handlePublish = async () => {
        setNotice("");
        setPublishResults([]);

        if (!session) {
            setNotice("Sign in before publishing.");
            return;
        }

        if (!draft.trim()) {
            setNotice("Write some content before publishing.");
            return;
        }

        if (selectedPlatforms.length === 0) {
            setNotice("Choose at least one connected platform.");
            return;
        }

        setPublishing(true);
        try {
            let response;

            if (attachedMedia) {
                const formData = new FormData();
                formData.append("content", draft);
                formData.append("platforms", JSON.stringify(selectedPlatforms));
                formData.append("media", attachedMedia);

                response = await fetch(`${API_BASE}/api/posts/publish`, {
                    method: "POST",
                    credentials: "include",
                    body: formData,
                });
            } else {
                response = await fetch(`${API_BASE}/api/posts/publish`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    credentials: "include",
                    body: JSON.stringify({
                        content: draft,
                        platforms: selectedPlatforms,
                    }),
                });
            }

            const data = await response.json();
            if (!response.ok && !data.results) {
                throw new Error(data.error || "Publishing failed");
            }

            const results = data.results || [];
            setPublishResults(results);

            const xFallback = results.find((result) => result.platform === "x" && result.fallback_url);
            if (xFallback?.fallback_url) {
                window.open(xFallback.fallback_url, "_blank", "noopener,noreferrer");
            }

            const successCount = results.filter((result) => result.success).length;
            if (successCount > 0) {
                const linkedInSuccess = results.some((result) => result.platform === "linkedin" && result.success);
                const xComposerOpened = results.some((result) => result.platform === "x" && result.fallback_url);

                if (linkedInSuccess && xComposerOpened) {
                    setNotice("LinkedIn was published. X composer opened with the same draft.");
                } else if (xComposerOpened) {
                    setNotice("X composer opened with your draft.");
                } else {
                    setNotice(successCount === selectedPlatforms.length ? "Post published." : "Post published partially.");
                }

                if (successCount === selectedPlatforms.length) {
                    setAttachedMedia(null);
                }
            } else {
                setNotice("Publishing failed on all selected platforms.");
            }
        } catch (error) {
            console.error(error);
            setNotice(error.message);
        } finally {
            setPublishing(false);
        }
    };

    const renderChannelPane = () => {
        if (activeChannel === "all") {
            return (
                <Hero
                    draft={draft}
                    setDraft={setDraft}
                    attachedMedia={attachedMedia}
                    selectedPlatforms={selectedPlatforms}
                    togglePlatform={togglePlatform}
                    platforms={session.platforms}
                    notice={notice}
                    publishResults={publishResults}
                    onPublish={handlePublish}
                    onPickMedia={handlePickMedia}
                    onRemoveMedia={handleRemoveMedia}
                    publishing={publishing}
                    xLimit={xLimit}
                    user={session.user}
                />
            );
        }

        const channelTitle = activeChannel === "x" ? "X" : "LinkedIn";

        return (
            <div className="min-h-screen bg-[linear-gradient(180deg,_rgba(12,10,9,0.96)_0%,_rgba(17,24,39,0.92)_50%,_rgba(2,6,23,1)_100%)] px-6 py-8 lg:px-10">
                <div className="mx-auto max-w-6xl">
                    <section className="rounded-[2rem] border border-white/10 bg-white/5 p-8 shadow-2xl shadow-black/20 backdrop-blur">
                        <p className="text-sm uppercase tracking-[0.3em] text-stone-400">Channel View</p>
                        <h1 className="mt-4 text-4xl font-black text-white">{channelTitle}</h1>
                        <div className="mt-8 rounded-[1.75rem] border border-dashed border-white/10 bg-slate-950/70 px-8 py-24 text-center">
                            <p className="text-lg font-semibold text-white">This channel page is empty for now.</p>
                            <p className="mt-3 text-sm text-stone-400">Select All Channels to return to the full composer and publish view.</p>
                        </div>
                    </section>
                </div>
            </div>
        );
    };

    if (loading) {
        return <div className="min-h-screen bg-stone-950 text-stone-100 flex items-center justify-center">Loading LinkNest...</div>;
    }

    if (!session) {
        return (
            <div className="min-h-screen bg-[radial-gradient(circle_at_top,_rgba(59,130,246,0.18),_transparent_32%),linear-gradient(180deg,_#0c0a09_0%,_#111827_48%,_#020617_100%)] text-stone-100 px-6 py-10">
                <div className="mx-auto flex min-h-[calc(100vh-5rem)] max-w-6xl flex-col justify-center gap-8 lg:flex-row lg:items-center">
                    <div className="max-w-2xl space-y-6">
                        <p className="text-sm uppercase tracking-[0.35em] text-sky-300">LinkNest MVP</p>
                        <h1 className="text-5xl font-black leading-tight text-white sm:text-6xl">
                            Write once. Publish to X and LinkedIn from one clean dashboard.
                        </h1>
                        <p className="max-w-xl text-lg text-stone-300">
                            Sign in with Google, connect your accounts, and ship text posts instantly. No scheduling, no clutter, no paid wall for the MVP.
                        </p>
                        {notice && (
                            <div className="inline-flex rounded-full border border-sky-400/30 bg-sky-400/10 px-4 py-2 text-sm text-sky-100">
                                {notice}
                            </div>
                        )}
                        <div className="flex flex-wrap gap-4">
                            <button
                                type="button"
                                onClick={handleGoogleLogin}
                                className="rounded-full bg-sky-400 px-6 py-3 text-sm font-semibold text-slate-950 transition hover:bg-sky-300"
                            >
                                Continue with Google
                            </button>
                            <a
                                href="/Loginpage"
                                className="rounded-full border border-stone-700 px-6 py-3 text-sm font-semibold text-stone-100 transition hover:border-stone-500 hover:bg-white/5"
                            >
                                Use email login
                            </a>
                        </div>
                    </div>

                    <div className="w-full max-w-xl rounded-[2rem] border border-white/10 bg-white/5 p-6 shadow-2xl shadow-sky-950/20 backdrop-blur">
                        <div className="rounded-[1.5rem] border border-white/10 bg-slate-950/80 p-6">
                            <div className="mb-6 flex items-center justify-between">
                                <div>
                                    <p className="text-sm text-stone-400">Launch flow</p>
                                    <h2 className="text-2xl font-bold text-white">Creator workspace</h2>
                                </div>
                                <span className="rounded-full border border-emerald-400/30 bg-emerald-400/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.25em] text-emerald-200">
                                    Free MVP
                                </span>
                            </div>
                            <div className="space-y-4 text-sm text-stone-300">
                                <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                                    1. Sign in with Google.
                                </div>
                                <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                                    2. Connect your X and LinkedIn personal profile.
                                </div>
                                <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                                    3. Publish one text post to either platform or both.
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-stone-950 text-stone-100">
            <Sidebar
                user={session.user}
                platforms={session.platforms}
                activeChannel={activeChannel}
                busyPlatform={busyPlatform}
                onConnect={handleConnectPlatform}
                onDisconnect={handleDisconnectPlatform}
                onOpenChannel={handleOpenChannel}
                onLogout={handleLogout}
            />
            <main className="lg:ml-80">
                {renderChannelPane()}
            </main>
        </div>
    );
}
