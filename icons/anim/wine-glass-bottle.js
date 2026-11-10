import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector("g");
  animate(bottle, { rotate: [0, -12, 3, 0] }, { duration: τmid });

  const wine = el.querySelector('path[d="M3 13h8"]');
  animate(wine, { rotate: [0, 12, -8, 0] }, { duration: τmid, delay: τmin });
}
