import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("rect");
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const cells = ["M6", "M10", "M14"].map((x) => el.querySelector(`path[d="${x} 10v4"]`));
  animate(cells, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
  const tip = el.querySelector('path[d="M22 14v-4"]');
  animate(tip, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
