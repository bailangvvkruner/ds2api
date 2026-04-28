import { Activity, Zap, TrendingUp, Clock, Timer } from 'lucide-react'

export default function StatisticsPanel({ queueStatus, t }) {
    const formatNumber = (num) => {
        if (num === undefined || num === null) return '0'
        if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M'
        if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
        return num.toString()
    }

    const formatTime = (seconds) => {
        if (seconds === undefined || seconds === null || seconds === 0) return '0'
        if (seconds < 1) return (seconds * 1000).toFixed(0) + 'ms'
        return seconds.toFixed(2) + 's'
    }

    const stats = queueStatus || {}

    const cards = [
        {
            icon: Activity,
            label: t('statistics.todayRequests'),
            value: formatNumber(stats.today_requests),
            subLabel: t('statistics.requestsUnit'),
            color: 'text-blue-500',
            bgColor: 'bg-blue-500/10',
        },
        {
            icon: Zap,
            label: t('statistics.todayTokens'),
            value: `${formatNumber(stats.today_input_tokens)} / ${formatNumber(stats.today_output_tokens)}`,
            subLabel: t('statistics.tokenDirection'),
            color: 'text-amber-500',
            bgColor: 'bg-amber-500/10',
        },
        {
            icon: TrendingUp,
            label: t('statistics.totalTokens'),
            value: `${formatNumber(stats.total_input_tokens)} / ${formatNumber(stats.total_output_tokens)}`,
            subLabel: t('statistics.tokenDirection'),
            color: 'text-emerald-500',
            bgColor: 'bg-emerald-500/10',
        },
        {
            icon: Clock,
            label: t('statistics.performance'),
            value: `${formatNumber(stats.rpm)} / ${formatNumber(stats.tpm)}`,
            subLabel: t('statistics.rpmTpm'),
            color: 'text-purple-500',
            bgColor: 'bg-purple-500/10',
        },
        {
            icon: Timer,
            label: t('statistics.avgResponseTime'),
            value: formatTime(stats.avg_response_time),
            subLabel: t('statistics.avgResponseUnit'),
            color: 'text-cyan-500',
            bgColor: 'bg-cyan-500/10',
        },
    ]

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
            {cards.map((card, i) => (
                <div key={i} className="bg-card border border-border rounded-xl p-4 flex flex-col justify-between shadow-sm relative overflow-hidden group">
                    <div className={`absolute right-0 top-0 p-4 opacity-5 group-hover:opacity-10 transition-opacity ${card.bgColor}`}>
                        <card.icon className="w-16 h-16" />
                    </div>
                    <p className="text-xs font-medium text-muted-foreground uppercase tracking-widest">{card.label}</p>
                    <div className="mt-2 flex items-baseline gap-1">
                        <span className={`text-2xl font-bold ${card.color}`}>{card.value}</span>
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">{card.subLabel}</p>
                </div>
            ))}
        </div>
    )
}
