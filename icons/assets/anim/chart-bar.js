import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = ["M7 16h8", "M7 11h12", "M7 6h3"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
