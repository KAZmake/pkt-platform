export interface FarmProperties {
  id: string;
  name: string;
  borrower_id?: string;
  activity_type: 'crop_farming' | 'livestock' | 'mixed';
  district: string;
  area_ha?: number;
  /** Only for employee+ */
  loan_status?: 'active' | 'overdue' | 'closed';
  loan_amount?: number;
}

export interface DistrictProperties {
  id: string;
  name: string;
  farms_count?: number;
  portfolio_sum?: number;
}

export type MapLayerGroup =
  | 'districts'
  | 'farms'
  | 'activity'
  | 'coverage'
  | 'portfolio'
  | 'loan-status';

export interface LayerToggle {
  id: MapLayerGroup;
  label: string;
  icon: string;
  employeeOnly: boolean;
  enabled: boolean;
}
