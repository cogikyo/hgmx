import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const door = el.querySelector("g[transform-origin]");
  animate(door, { scaleX: [1, 0.8, 1] }, { duration: τmid });

  const gear = el.querySelector("g.gear");
  animate(gear, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmid, delay: τmin });
}
