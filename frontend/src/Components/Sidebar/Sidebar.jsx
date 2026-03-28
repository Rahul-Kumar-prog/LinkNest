import linknestLogo from "../../assets/linknest-logo.svg";

const XLogo = ({ className = "h-3 w-3" }) => (
    <svg viewBox="0 0 24 24" className={className} fill="currentColor" aria-hidden="true">
        <path d="M18.901 2H22l-6.768 7.736L23.194 22H16.96l-4.884-7.443L5.56 22H2.458l7.237-8.272L1.806 2H8.2l4.415 6.73L18.901 2Zm-1.087 18.112h1.717L7.267 3.79H5.425l12.389 16.322Z" />
    </svg>
);

const ChannelAvatar = ({ type, imageUrl, accent, badge }) => {
    if (type === "all") {
        return (
            <div className="flex h-12 w-12 items-center justify-center rounded-full border border-white/30 bg-transparent text-white">
                <svg viewBox="0 0 24 24" className="h-6 w-6" fill="none" stroke="currentColor" strokeWidth="1.7">
                    <circle cx="8" cy="8" r="2.25" />
                    <circle cx="16" cy="8" r="2.25" />
                    <circle cx="8" cy="16" r="2.25" />
                    <circle cx="16" cy="16" r="2.25" />
                </svg>
            </div>
        );
    }

    if (imageUrl) {
        return (
            <div className="relative h-12 w-12 shrink-0">
                <img
                    src={imageUrl}
                    alt={type}
                    className="h-12 w-12 rounded-full object-cover ring-1 ring-white/10"
                />
                {badge ? (
                    <span className={["absolute -bottom-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full border border-stone-950 text-[10px] font-bold text-white", badge].join(" ")}>
                        {type === "linkedin" ? "in" : <XLogo className="h-2.5 w-2.5" />}
                    </span>
                ) : null}
            </div>
        );
    }

    return (
        <div className="relative h-12 w-12 shrink-0">
            <div className={["flex h-12 w-12 items-center justify-center rounded-full text-sm font-bold text-white", accent].join(" ")}>
                {type === "linkedin" ? "in" : <XLogo className="h-2.5 w-2.5" />}
            </div>
            {badge ? (
                <span className={["absolute -bottom-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full border border-stone-950 text-[10px] font-bold text-white", badge].join(" ")}>
                    {type === "linkedin" ? "in" : <XLogo className="h-2.5 w-2.5" />}
                </span>
            ) : null}
        </div>
    );
};

const channels = [
    { id: "all", name: "All Channels", accent: "bg-stone-700" },
    { id: "linkedin", name: "LinkedIn", accent: "bg-blue-700", badge: "bg-blue-600" },
    { id: "x", name: "X", accent: "bg-white text-slate-950", badge: "bg-white text-slate-950" },
];

export default function Sidebar({ user, platforms, activeChannel, onOpenChannel, onLogout }) {
    return (
        <aside className="border-b border-white/10 bg-slate-950/95 px-6 py-6 lg:fixed lg:left-0 lg:top-0 lg:h-screen lg:w-80 lg:border-b-0 lg:border-r">
            <div className="flex h-full flex-col">
                <div>
                    <img
                        src={linknestLogo}
                        alt="LinkNest"
                        className="h-8 w-auto"
                    />
                    <h1 className="mt-4 text-3xl font-black text-white">Creator dashboard</h1>
                    <p className="mt-3 text-sm text-stone-400">
                        Connect your profiles, write one text post, and publish instantly.
                    </p>
                </div>

                <div className="mt-8 rounded-[2rem] border border-white/10 bg-white/[0.04] p-5">
                    <div className="flex items-center gap-4">
                        {user?.picture ? (
                            <img
                                src={user.picture}
                                alt={user?.name || "User"}
                                className="h-12 w-12 rounded-2xl object-cover"
                            />
                        ) : (
                            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-fuchsia-500 text-sm font-bold text-white">
                                {(user?.name || user?.email || "U").slice(0, 2).toUpperCase()}
                            </div>
                        )}
                        <div className="min-w-0">
                            <p className="truncate text-base font-semibold text-white">{user?.name || "LinkNest user"}</p>
                            <p className="truncate text-sm text-stone-400">{user?.email}</p>
                        </div>
                    </div>
                </div>

                <div className="mt-8 flex-1">
                    <div className="flex items-center justify-between">
                        <h2 className="text-sm font-semibold uppercase tracking-[0.3em] text-stone-500">Channels</h2>
                        <span className="rounded-full border border-white/10 px-3 py-1 text-xs text-stone-300">3 listed</span>
                    </div>

                    <div className="mt-4 space-y-2">
                        {channels.map((channel) => {
                            const connection = platforms?.[channel.id];
                            const connected = channel.id === "all" ? true : !!connection?.connected;
                            const label = channel.id === "all"
                                ? channel.name
                                : channel.id === "x"
                                    ? (connection?.name || connection?.username || connection?.user_id || channel.name)
                                    : (connection?.name || channel.name);
                            const rowClass = activeChannel === channel.id
                                ? "w-full rounded-[0.9rem] bg-indigo-900/80 px-3 py-3 text-left transition hover:bg-indigo-800/80"
                                : "w-full rounded-[0.9rem] bg-transparent px-3 py-3 text-left transition hover:bg-white/[0.04]";

                            return (
                                <button
                                    key={channel.id}
                                    type="button"
                                    onClick={() => onOpenChannel(channel.id)}
                                    className={rowClass}
                                >
                                    <div className="flex items-center gap-3">
                                        <ChannelAvatar
                                            type={channel.id}
                                            imageUrl={connection?.profile_picture}
                                            accent={channel.accent}
                                            badge={channel.badge}
                                        />
                                        <div className="min-w-0 flex-1">
                                            <p className="truncate text-sm font-semibold text-white">{label}</p>
                                            {channel.id !== "all" && !connected ? (
                                                <p className="truncate text-xs text-stone-500">Not connected</p>
                                            ) : null}
                                        </div>
                                        <span className="text-sm font-semibold text-stone-200">0</span>
                                    </div>
                                </button>
                            );
                        })}
                    </div>
                </div>

                <button
                    type="button"
                    onClick={onLogout}
                    className="mt-8 rounded-full border border-white/10 px-5 py-3 text-sm font-semibold text-stone-100 transition hover:bg-white/[0.06]"
                >
                    Logout
                </button>
            </div>
        </aside>
    );
}
