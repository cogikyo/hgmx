import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = ["M8 18v-2", "M12 18v-4", "M16 18v-6"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
