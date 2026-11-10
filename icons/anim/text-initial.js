import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const initial = el.querySelectorAll('path[d^="m3 12"], path[d^="M3.92"]');
  animate(initial, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const lines = ["M15 5h6", "M15 12h6", "M3 19h18"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { translateX: [4, -1, 0] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
