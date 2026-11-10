import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const [lens, mark] = el.querySelectorAll("g");
  animate(lens, { scale: [1.35, 1] }, { duration: τmax });
  animate(lens, { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });

  animate(mark, { rotate: [0, -18, 14, -8, 0] }, { duration: τmid, delay: τmin });

  const dot = el.querySelector('path[d="M11 15h.01"]');
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
