import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector('path[d^="M21.25"]');
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const hands = el.querySelector('path[d^="M12 6v6"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmax });
  const mark = el.querySelector("g");
  animate(mark, { rotate: [0, -18, 14, -8, 0] }, { duration: τmid, delay: τmin });
  const dot = el.querySelector('path[d^="M20 21"]');
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
