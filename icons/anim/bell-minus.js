import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const bell = el.querySelector("g");
  animate(bell, { rotate: [0, -18, 14, -8, 0] }, { duration: τmax });
  const clapper = el.querySelector('path[d^="M10.268"]');
  animate(clapper, { translateX: [0, 2, -2, 0] }, { duration: τmid, delay: τmin });
  const minus = el.querySelector('path[d="M15 8h6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
