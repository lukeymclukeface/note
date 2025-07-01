export default function CliSetupLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // This layout excludes the Navigation component
  // to provide a clean setup experience
  return <>{children}</>;
}
