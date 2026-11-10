import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('path[d^="M18 8"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const phone = el.querySelector("rect");
  animate(phone, { rotate: [0, -10, 10, -6, 0] }, { duration: τmid, delay: τmin });
}
