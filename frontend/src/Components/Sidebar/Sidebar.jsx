import { useState } from 'react';

const DropdownIcon = () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
    </svg>
);

const Avatar = ({ initials, bgColor }) => (
    <div className={`w-8 h-8 ${bgColor} rounded-full flex items-center justify-center text-white text-xs font-medium`}>
        {initials}
    </div>
);

export default function Sidebar() {
    const [showMore, setShowMore] = useState(false);

    const userChannels = [
        { id: 'all', name: 'All Channels', count: 2, isHeader: true },
        { id: 'twitter', name: 'Twitter / X', avatar: 'TW', avatarColor: 'bg-blue-400', count: '12' },
        { id: 'linkedin', name: 'LinkedIn', avatar: 'LI', avatarColor: 'bg-blue-700', count: '8' },
    ];

    const moreChannels = [
        { id: 'instagram', name: 'Connect Instagram', disabled: true },
        { id: 'threads', name: 'Connect Threads', disabled: true },
        { id: 'bluesky', name: 'Connect Bluesky', disabled: true },
        { id: 'facebook', name: 'Connect Facebook', disabled: true },
        { id: 'tiktok', name: 'Connect TikTok', disabled: true },
        { id: 'youtube', name: 'Connect YouTube', disabled: true },
    ];

    const handlePlatform = async (channelId) => {
        // const actions = {
        //     twitter: () => console.log("you are on Twitter"),
        //     linkedin: () => console.log("you are on LinkedIn"),
        // };
        try {
            const response = await fetch(`http://localhost:8080/api/${channelId}`);
            if (!response.ok) {
            const text = await response.text();
            throw new Error(`HTTP ${response.status}: ${text}`);
        }
            const data = await response.json();
            console.log(data);
        }catch(error) {
            console.error("API error:", error);
        }

        // actions[channelId]?.();
    };

    return (
        <aside className="fixed left-0 top-0 h-screen w-64 bg-gray-900 border-r border-gray-800 flex flex-col">
            {/* Header */}
            <div className="p-4 border-b border-gray-800">
                <h1 className="text-2xl font-bold text-white">LinkNest</h1>
            </div>

            {/* Channels Section */}
            <div className="flex-1 overflow-y-auto">
                <div className="p-4">
                    <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
                        Channels
                    </h2>

                    <div className="space-y-1">
                        {/* All Channels */}
                        {userChannels.slice(0, 1).map((channel) => (
                            <button
                                key={channel.id}
                                className={`w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm transition-colors ${channel.isHeader
                                    ? 'bg-gray-800 text-white'
                                    : 'hover:bg-gray-800 text-gray-300'
                                    }`}
                            >
                                <span className="font-medium">{channel.name}</span>
                                <span className="bg-gray-700 text-gray-300 text-xs px-2 py-0.5 rounded-full">
                                    {channel.count}
                                </span>
                            </button>
                        ))}

                        {/* User Channels */}
                        {userChannels.slice(1).map((channel) => (
                            <button
                                type='button'
                                key={channel.id}
                                onClick={() => handlePlatform(channel.id)}
                                className="cursor-pointer w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm hover:bg-gray-800 transition-colors group"
                            >
                                <div className="flex items-center gap-3">
                                    <Avatar initials={channel.avatar} bgColor={channel.avatarColor} />
                                    <span className="text-gray-300 group-hover:text-white">
                                        {channel.name}
                                    </span>
                                </div>
                            </button>
                        ))}

                        {/* Show More Channels */}
                        <button
                            onClick={() => setShowMore(!showMore)}
                            className="w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm text-gray-500 hover:text-gray-300 hover:bg-gray-800 transition-colors"
                        >
                            <span>Show more channels</span>
                            <div
                                className={`transform transition-transform duration-200 ${showMore ? 'rotate-180' : ''
                                    }`}
                            >
                                <DropdownIcon />
                            </div>
                        </button>

                        {/* Expanded Channels */}
                        {showMore && (
                            <div className="space-y-1 mt-1 animate-fadeIn">
                                {moreChannels.map((channel) => (
                                    <button
                                        key={channel.id}
                                        disabled
                                        className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-600 cursor-not-allowed"
                                    >
                                        <div className="w-8 h-8 bg-gray-800 rounded-full flex items-center justify-center">
                                            <div className="w-4 h-4 bg-gray-700 rounded" />
                                        </div>
                                        <span>{channel.name}</span>
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            </div>

            {/* Footer */}
            <div className="p-4 border-t border-gray-800">
                <div className="flex items-center gap-3 px-3 py-2">
                    <Avatar initials="U" bgColor="bg-purple-500" />
                    <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-white truncate">User Name</p>
                        <p className="text-xs text-gray-500 truncate">user@example.com</p>
                    </div>
                </div>
            </div>
        </aside>
    );
}

