import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = ["M7 6h3", "M7 11h8", "M7 16h12"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { scaleX: [1, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
