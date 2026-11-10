import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const die = el.querySelector("g");
  animate(die, { translateY: [0, -3, 1, 0], rotate: [0, -15, 6, 0] }, { duration: τmax });
  const pips = (ds) => ds.map((d) => el.querySelector(`path[d="${d}h.01"]`));
  animate(pips(["M8 12", "M16 12"]), { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(pips(["M8 8", "M16 8"]), { translateY: [4, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(pips(["M8 16", "M16 16"]), { translateY: [-4, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
