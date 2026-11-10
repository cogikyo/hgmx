import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const cart = el.querySelector("g");
  animate(cart, { translateX: [0, 4, -2, 0], rotate: [0, -8, 4, 0] }, { duration: τmax });

  const minus = el.querySelector('path[d="M16 5h6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
