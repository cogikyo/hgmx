import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector("rect");
  animate(walls, { scaleY: [0, 1.1, 1] }, { duration: τmid });

  const windows = ["14", "10", "6"].flatMap((y) => ["8", "12", "16"].map((x) => el.querySelector(`path[d="M${x} ${y}h.01"]`)));
  animate(windows, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });

  const door = el.querySelector('path[d^="M9 22"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
