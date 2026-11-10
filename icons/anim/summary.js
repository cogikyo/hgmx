import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = ["M15 4H7", "M7 9h12", "M7 14h7"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });
  const route = el.querySelector('path[d^="M3 4v13"]');
  animate(route, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
