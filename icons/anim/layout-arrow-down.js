import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const tiles = el.querySelectorAll("rect");
  animate(tiles, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });

  const arrow = el.querySelectorAll("path");
  animate(arrow, { translateY: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
