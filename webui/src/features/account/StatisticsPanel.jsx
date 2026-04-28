import { useEffect, useState } from 'react'
import { Activity, Zap, Clock, TrendingUp } from 'lucide-react'

export default function StatisticsPanel({ apiFetch, t }) {
    const [stats, setStats] = useState(null)
    const [loading, setLoading] = useState(true)

    const formatNumber = (num) => {
        if (num === undefined || num === null) return '0'
        if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M'
        if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
        return num.toString()
    }

    useEffect(() => {
        const fetchStats = async () => {
            try {
                const res = await apiFetch('/admin/statistics')
                if (res.ok) {
                    const data = await res.json()
                    setStats(data)
                }
            } catch (e) {
                console.error('Failed to fetch statistics:', e)
            } finally {
                setLoading(false)
            }
        }

        fetchStats()
        const interval = setInterval(fetchStats, 10000)
        return () => clearInterval(interval)
    }, [apiFetch])

    if (loading) {
        return (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {[1, 2, 3, 4].map(i => (
                    <div key={i} className="bg-card border border-border rounded-xl p-4 animate-pulse">
                        <div className="h-4 bg-muted rounded w-24 mb-2"></div>
                        <div className="h-8 bg-muted rounded w-16"></div>
                    </div>
                ))}
            </div>
        )
    }

    if (!stats) {
        return (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {[1, 2, 3, 4].map(i => (
                    <div key={i} className="bg-card border border-border rounded-xl p-4">
                        <p className="text-xs text-muted-foreground">{i === 0 ? t('statistics.todayRequests') : i === 1 ? t('statistics.todayTokens') : i === 2 ? t('statistics.totalTokens') : t('statistics.performance')}</p>
                        <p className="text-2xl font-bold text-muted-foreground mt-2">-</p>
                    </div>
                ))}
            </div>
        )
    }

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
    ]

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
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
