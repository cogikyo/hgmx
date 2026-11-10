import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g[transform-origin]");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });

  const bars = ["M8 7v10", "M12 7v10", "M17 7v10"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
