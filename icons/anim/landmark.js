import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ground = el.querySelector('path[d="M3 22h18"]');
  animate(ground, { scaleX: [0, 1] }, { duration: τmid });

  const columns = ["M6", "M10", "M14", "M18"].map((x) => el.querySelector(`path[d="${x} 18v-7"]`));
  animate(columns, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const roof = el.querySelector('path[d^="M11.119"]');
  animate(roof, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
