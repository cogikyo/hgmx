import { animate, stagger } from "../motion.js";

export default function (el) {
  const tiles = el.querySelectorAll("rect");
  animate(tiles, { scale: [0, 1.25, 1], rotate: [-20, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
