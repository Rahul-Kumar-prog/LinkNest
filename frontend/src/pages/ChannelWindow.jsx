import { useMemo } from "react";
import { useParams } from "react-router-dom";

const channelMeta = {
    x: { title: "X", accent: "from-sky-500/30 to-cyan-400/10" },
    linkedin: { title: "LinkedIn", accent: "from-blue-600/30 to-indigo-400/10" },
};

export default function ChannelWindow() {
    const { channelId } = useParams();
    const channel = useMemo(() => channelMeta[channelId] || { title: "Channel", accent: "from-stone-500/20 to-stone-400/10" }, [channelId]);

    return (
        <div className="min-h-screen bg-[linear-gradient(180deg,_rgba(12,10,9,0.96)_0%,_rgba(17,24,39,0.92)_50%,_rgba(2,6,23,1)_100%)] px-6 py-8 text-stone-100 lg:px-10">
            <div className="mx-auto max-w-5xl space-y-8">
                <header className={`rounded-[2rem] border border-white/10 bg-gradient-to-br ${channel.accent} p-8 shadow-2xl shadow-black/30 backdrop-blur`}>
                    <p className="text-sm uppercase tracking-[0.3em] text-stone-300">Dedicated Window</p>
                    <h1 className="mt-4 text-4xl font-black text-white">{channel.title}</h1>
                    <p className="mt-3 max-w-2xl text-base text-stone-200">
                        This channel opened in its own window. You can build a dedicated {channel.title} workflow here next.
                    </p>
                </header>

                <section className="rounded-[2rem] border border-dashed border-white/10 bg-slate-950/70 p-10 text-center shadow-2xl shadow-black/30">
                    <p className="text-sm uppercase tracking-[0.3em] text-stone-500">Placeholder</p>
                    <h2 className="mt-4 text-2xl font-bold text-white">Separate {channel.title} window is ready</h2>
                    <p className="mt-3 text-sm text-stone-400">
                        All Channels continues to show the main composer in the original dashboard window.
                    </p>
                </section>
            </div>
        </div>
    );
}
