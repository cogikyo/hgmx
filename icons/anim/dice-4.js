import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const die = el.querySelector("g");
  animate(die, { translateY: [0, -3, 1, 0], rotate: [0, -15, 6, 0] }, { duration: τmax });
  const pip = (d) => el.querySelector(`path[d="${d}h.01"]`);
  animate([pip("M16 8"), pip("M8 16")], { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: τmin, delay: τmin });
  animate([pip("M8 8"), pip("M16 16")], { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
