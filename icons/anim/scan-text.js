import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g[transform-origin]");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });

  const lines = ["M7 8h8", "M7 12h10", "M7 16h6"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
