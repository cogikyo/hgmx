import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const cart = el.querySelector("g");
  animate(cart, { translateX: [0, 4, -2, 0], rotate: [0, -8, 4, 0] }, { duration: τmax });

  const plus = el.querySelectorAll('path[d="M16 5h6"], path[d="M19 2v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
