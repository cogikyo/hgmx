import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const feather = el.querySelector("g");
  animate(feather, { rotate: [0, -12, 8, 0] }, { duration: τmax });

  const plus = el.querySelectorAll('path[d="M2 7h6"], path[d="M5 4v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
