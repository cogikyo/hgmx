import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d^="m7 7"]');
  animate(left, { translateX: [0, -1, 0.5, 0] }, { duration: τmax });

  const right = el.querySelector('path[d^="m17 7"]');
  animate(right, { translateX: [0, 1, -0.5, 0] }, { duration: τmax });

  const dots = ["8", "12", "16"].map((x) => el.querySelector(`path[d="M${x} 12h.01"]`));
  animate(dots, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
