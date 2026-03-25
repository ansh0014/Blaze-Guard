import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { DashboardLayout } from './src/modules/authority/views/DashboardLayout';
import { OperationalMap } from './src/modules/authority/views/OperationalMap';
import { LandingPage } from './src/LandingPage';
import { AlertsPanel } from './src/modules/authority/views/AlertsPanel';
import { ResourcesView } from './src/modules/authority/views/ResourcesView';
import { RiskAnalyticsPanel } from './src/modules/authority/views/RiskAnalyticsPanel';
import { AIReportView } from './src/modules/authority/views/AIReportView';
import { FeedbackInterface } from './src/modules/authority/views/FeedbackInterface';
import { SafetyPortal } from './src/modules/citizen/views/SafetyPortal';
import { ProtectedRoute } from './src/components/ProtectedRoute';
import { LoginPage } from './src/components/LoginPage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<LoginPage />} />

        {/* Authority Routes - Protected */}
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <DashboardLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="/dashboard/map" replace />} />
          <Route path="map" element={<OperationalMap />} />
          <Route path="alerts" element={<AlertsPanel />} />
          <Route path="resources" element={<ResourcesView />} />
          <Route path="analytics" element={<RiskAnalyticsPanel />} />
          <Route path="report" element={<AIReportView />} />
          <Route path="feedback" element={<FeedbackInterface />} />
        </Route>

        {/* Citizen Routes */}
        <Route path="/citizen" element={<SafetyPortal />} />

        {/* Redirects for legacy/broken links */}
        <Route path="/authority/map" element={<Navigate to="/dashboard/map" replace />} />
        <Route path="*" element={<Navigate to="/dashboard/map" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
