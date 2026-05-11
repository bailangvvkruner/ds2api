import { Coins, Gauge, Send, Timer, TrendingUp } from 'lucide-react'

export default function StatisticsPanel({ queueStatus, t }) {
    if (!queueStatus) {
        return null
    }

    const formatNumber = (value) => {
        const num = Number(value || 0)
        if (num >= 1000000000) return `${(num / 1000000000).toFixed(2)}B`
        if (num >= 1000000) return `${(num / 1000000).toFixed(2)}M`
        if (num >= 1000) return `${(num / 1000).toFixed(1)}K`
        return new Intl.NumberFormat().format(num)
    }

    const formatDuration = (ms) => {
        const num = Number(ms || 0)
        if (num >= 1000) return `${(num / 1000).toFixed(2)}s`
        return `${Math.round(num)}ms`
    }

    const todayTokens = Number(queueStatus.today_input_tokens || 0) + Number(queueStatus.today_output_tokens || 0)
    const totalTokens = Number(queueStatus.total_input_tokens || 0) + Number(queueStatus.total_output_tokens || 0)

    const cards = [
        {
            icon: Send,
            label: t('statistics.todayRequests'),
            value: formatNumber(queueStatus.today_requests),
            detail: t('statistics.requestsUnit'),
            color: 'text-blue-500',
            border: 'border-blue-500/30',
            bg: 'from-blue-500/15 to-blue-500/5',
        },
        {
            icon: Coins,
            label: t('statistics.todayTokens'),
            value: formatNumber(todayTokens),
            detail: `${t('statistics.inputOutput')}: ${formatNumber(queueStatus.today_input_tokens)} / ${formatNumber(queueStatus.today_output_tokens)}`,
            color: 'text-amber-500',
            border: 'border-amber-500/30',
            bg: 'from-amber-500/15 to-amber-500/5',
        },
        {
            icon: TrendingUp,
            label: t('statistics.totalTokens'),
            value: formatNumber(totalTokens),
            detail: `${t('statistics.inputOutput')}: ${formatNumber(queueStatus.total_input_tokens)} / ${formatNumber(queueStatus.total_output_tokens)}`,
            color: 'text-emerald-500',
            border: 'border-emerald-500/30',
            bg: 'from-emerald-500/15 to-emerald-500/5',
        },
        {
            icon: Gauge,
            label: t('statistics.performance'),
            value: `${formatNumber(queueStatus.rpm)} / ${formatNumber(queueStatus.tpm)}`,
            detail: 'RPM / TPM',
            color: 'text-purple-500',
            border: 'border-purple-500/30',
            bg: 'from-purple-500/15 to-purple-500/5',
        },
        {
            icon: Timer,
            label: t('statistics.average'),
            value: `${formatDuration(queueStatus.average_response_ms)} / ${formatDuration(queueStatus.average_time_ms)}`,
            detail: t('statistics.responseTime'),
            color: 'text-cyan-500',
            border: 'border-cyan-500/30',
            bg: 'from-cyan-500/15 to-cyan-500/5',
        },
    ]

    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-5 gap-4">
            {cards.map((card) => (
                <div key={card.label} className={`relative bg-gradient-to-br ${card.bg} border ${card.border} rounded-xl p-4 shadow-sm overflow-hidden`}>
                    <div className="absolute -right-4 -top-4 opacity-10">
                        <card.icon className="w-20 h-20" />
                    </div>
                    <div className="relative space-y-2">
                        <div className="flex items-center gap-2">
                            <card.icon className={`w-4 h-4 ${card.color}`} />
                            <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider">{card.label}</p>
                        </div>
                        <div className={`text-2xl font-bold ${card.color}`}>{card.value}</div>
                        <p className="text-xs text-muted-foreground">{card.detail}</p>
                    </div>
                </div>
            ))}
        </div>
    )
}
