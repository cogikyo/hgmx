import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const phone = el.querySelector("g");
  animate(phone, { rotate: [0, -10, 10, -6, 0] }, { duration: τmax });

  const button = el.querySelector('path[d="M8 18h.01"]');
  animate(button, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
