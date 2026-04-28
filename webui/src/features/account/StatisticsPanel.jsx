import { Activity, Zap, TrendingUp, Clock, Timer, DollarSign, Send, Coins, Gauge } from 'lucide-react'

export default function StatisticsPanel({ queueStatus, t }) {
    const formatNumber = (num) => {
        if (num === undefined || num === null) return '0'
        if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M'
        if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
        return num.toLocaleString()
    }

    const formatTime = (seconds) => {
        if (seconds === undefined || seconds === null || seconds === 0) return '0'
        if (seconds < 1) return (seconds * 1000).toFixed(0) + 'ms'
        return seconds.toFixed(2) + 's'
    }

    const stats = queueStatus || {}
    const totalTokens = (stats.today_input_tokens || 0) + (stats.today_output_tokens || 0)

    const cards = [
        {
            icon: Send,
            label: t('statistics.todayRequests'),
            value: formatNumber(stats.today_requests),
            subLabel: t('statistics.requestsUnit'),
            gradient: 'from-blue-500/20 to-blue-600/10',
            border: 'border-blue-500/30',
            iconColor: 'text-blue-500',
        },
        {
            icon: Coins,
            label: t('statistics.todayTokens'),
            value: formatNumber(totalTokens),
            subLabel: t('statistics.tokenDirection2'),
            gradient: 'from-amber-500/20 to-amber-600/10',
            border: 'border-amber-500/30',
            iconColor: 'text-amber-500',
            details: `${formatNumber(stats.today_input_tokens)} / ${formatNumber(stats.today_output_tokens)}`,
        },
        {
            icon: TrendingUp,
            label: t('statistics.totalTokens'),
            value: formatNumber((stats.total_input_tokens || 0) + (stats.total_output_tokens || 0)),
            subLabel: t('statistics.cumulative'),
            gradient: 'from-emerald-500/20 to-emerald-600/10',
            border: 'border-emerald-500/30',
            iconColor: 'text-emerald-500',
            details: `${formatNumber(stats.total_input_tokens)} / ${formatNumber(stats.total_output_tokens)}`,
        },
        {
            icon: Gauge,
            label: t('statistics.performance'),
            value: `${formatNumber(stats.rpm)} / ${formatNumber(stats.tpm)}`,
            subLabel: 'RPM / TPM',
            gradient: 'from-purple-500/20 to-purple-600/10',
            border: 'border-purple-500/30',
            iconColor: 'text-purple-500',
        },
        {
            icon: Timer,
            label: t('statistics.avgResponseTime'),
            value: formatTime(stats.avg_response_time),
            subLabel: t('statistics.avgResponseUnit'),
            gradient: 'from-cyan-500/20 to-cyan-600/10',
            border: 'border-cyan-500/30',
            iconColor: 'text-cyan-500',
        },
    ]

    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
            {cards.map((card, i) => (
                <div
                    key={i}
                    className={`relative bg-gradient-to-br ${card.gradient} border ${card.border} rounded-xl p-4 shadow-sm overflow-hidden group hover:shadow-md transition-shadow`}
                >
                    <div className="absolute -right-4 -top-4 w-20 h-20 opacity-10 group-hover:opacity-20 transition-opacity">
                        <card.icon className="w-full h-full" />
                    </div>

                    <div className="relative">
                        <div className="flex items-center gap-2 mb-2">
                            <card.icon className={`w-4 h-4 ${card.iconColor}`} />
                            <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                                {card.label}
                            </p>
                        </div>

                        <div className="space-y-1">
                            <span className={`text-2xl font-bold ${card.iconColor}`}>
                                {card.value}
                            </span>
                            {card.details && (
                                <p className="text-xs text-muted-foreground/80">
                                    {card.details}
                                </p>
                            )}
                        </div>

                        <p className="text-xs text-muted-foreground mt-2">{card.subLabel}</p>
                    </div>
                </div>
            ))}
        </div>
    )
}
