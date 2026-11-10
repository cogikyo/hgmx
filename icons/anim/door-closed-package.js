import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const door = el.querySelector("g[transform-origin]");
  animate(door, { scaleX: [1, 0.8, 1] }, { duration: τmid });

  const box = el.querySelector("g.package");
  animate(box, { translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
