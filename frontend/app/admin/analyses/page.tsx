'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { useAuthStore } from '@/store/auth-store';
import { AdminRoute } from '@/components/AdminRoute';
import { Analysis, AnalysisListResponse } from '@/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ArrowLeft, Calendar, TrendingUp } from 'lucide-react';

function AnalysisHistoryPage() {
  const router = useRouter();
  const { clearAuth } = useAuthStore();
  const [analyses, setAnalyses] = useState<Analysis[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadAnalyses();
  }, []);

  const loadAnalyses = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.listAnalyses();
      // Response might be wrapped in a data property or be direct
      const responseData = response?.data || response;
      setAnalyses(responseData?.analyses || []);
    } catch (err: any) {
      if (err.response?.status === 401 || err.response?.status === 403) {
        clearAuth();
        router.push('/login');
      } else {
        setError('Failed to load analyses. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const getSentimentBadgeVariant = (sentiment: string) => {
    switch (sentiment) {
      case 'positive':
        return 'default';
      case 'negative':
        return 'destructive';
      case 'mixed':
        return 'secondary';
      default:
        return 'secondary';
    }
  };

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'success':
        return 'default';
      case 'failed':
        return 'destructive';
      case 'processing':
        return 'secondary';
      default:
        return 'secondary';
    }
  };

  if (isLoading) {
    return (
      <AdminRoute>
        <div className="flex items-center justify-center min-h-screen">
          <div className="text-center">
            <h1 className="text-2xl font-bold mb-4">Loading...</h1>
          </div>
        </div>
      </AdminRoute>
    );
  }

  return (
    <AdminRoute>
      <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="flex items-center gap-4">
            <Button variant="outline" onClick={() => router.push('/admin')}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Dashboard
            </Button>
            <h1 className="text-3xl font-bold">Analysis History</h1>
          </div>

          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {analyses.length === 0 ? (
            <Card>
              <CardContent className="pt-6">
                <p className="text-muted-foreground text-center">No analyses available yet.</p>
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-4">
              {analyses.map((analysis) => (
                <Card
                  key={analysis.id}
                  className="cursor-pointer hover:shadow-md transition-shadow"
                  onClick={() => router.push(`/admin/analyses/${analysis.id}`)}
                >
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <CardTitle className="mb-2">
                          Analysis from {new Date(analysis.period_start).toLocaleDateString()} to{' '}
                          {new Date(analysis.period_end).toLocaleDateString()}
                        </CardTitle>
                        <CardDescription className="flex items-center gap-4 mt-2">
                          <span className="flex items-center gap-1">
                            <Calendar className="h-4 w-4" />
                            {new Date(analysis.created_at).toLocaleString()}
                          </span>
                          {analysis.completed_at && (
                            <span className="flex items-center gap-1">
                              <TrendingUp className="h-4 w-4" />
                              Completed: {new Date(analysis.completed_at).toLocaleString()}
                            </span>
                          )}
                        </CardDescription>
                      </div>
                      <div className="flex gap-2">
                        <Badge variant={getSentimentBadgeVariant(analysis.sentiment)}>
                          {analysis.sentiment}
                        </Badge>
                        <Badge variant={getStatusBadgeVariant(analysis.status)}>
                          {analysis.status}
                        </Badge>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div>
                        <p className="text-sm text-muted-foreground mb-1">Summary</p>
                        <p className="text-sm line-clamp-2">{analysis.overall_summary}</p>
                      </div>
                      <div className="flex gap-4 text-sm">
                        <div>
                          <span className="text-muted-foreground">Feedbacks: </span>
                          <span className="font-semibold">{analysis.feedback_count}</span>
                        </div>
                        {analysis.new_feedback_count !== null && analysis.new_feedback_count !== undefined && (
                          <div>
                            <span className="text-muted-foreground">New: </span>
                            <span className="font-semibold">{analysis.new_feedback_count}</span>
                          </div>
                        )}
                        <div>
                          <span className="text-muted-foreground">Model: </span>
                          <span className="font-semibold">{analysis.model}</span>
                        </div>
                        <div>
                          <span className="text-muted-foreground">Tokens: </span>
                          <span className="font-semibold">{analysis.tokens.toLocaleString()}</span>
                        </div>
                      </div>
                      {analysis.key_insights && analysis.key_insights.length > 0 && (
                        <div>
                          <p className="text-sm text-muted-foreground mb-1">Key Insights</p>
                          <ul className="list-disc list-inside text-sm space-y-1">
                            {analysis.key_insights.slice(0, 3).map((insight, i) => (
                              <li key={i} className="line-clamp-1">{insight}</li>
                            ))}
                            {analysis.key_insights.length > 3 && (
                              <li className="text-muted-foreground">
                                +{analysis.key_insights.length - 3} more
                              </li>
                            )}
                          </ul>
                        </div>
                      )}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </AdminRoute>
  );
}

export default function AnalysesPage() {
  return <AnalysisHistoryPage />;
}
