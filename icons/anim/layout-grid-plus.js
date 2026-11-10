import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const tiles = el.querySelectorAll("rect");
  animate(tiles, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const plus = el.querySelectorAll("path");
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
